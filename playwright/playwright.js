const fs = require('fs');
const path = require('path');

const {Page, Browser, BrowserContext, Locator} = require('playwright-core');
const {chromium} = require('playwright-extra');
const {defaultTo} = require('lodash');
const StealthPlugin = require('puppeteer-extra-plugin-stealth');

chromium.use(StealthPlugin()); 
chromium.plugins.setDependencyDefaults('stealth/evasions/webgl.vendor', {
  vendor: 'Bob',
  renderer: 'Alice' 
});

const sleep= async (timeout) => {    await new Promise((resolve) => {
      setTimeout(() => {
        resolve();
      }, timeout);
    });
  }

class Scrapper {
  screenshotID = 1;
  /**
   * @type {Page}
   */
  page;
  /**
   * @type {Browser}
   */
  browser;
  /**
   * @type {BrowserContext}
   */
  context;

  /**
   *
   * @param {{
  * downloadsPath: string,
  * headless: boolean,
  * }} param0
   */
  init = async ({downloadsPath, headless}) => {
    this.browser = await chromium.launch({
      headless: defaultTo(headless, false),
      args: [
        '--no-sandbox'
      ],
      downloadsPath
    });
    this.context = await this.browser.newContext();
    this.page = await this.context.newPage();
  };

  close = async () => {
    await this.browser.close();
  };

  writeCookies = async (filepath) => {
    const cookies = await this.context.cookies();
    await fs.promises.writeFile(filepath, JSON.stringify(cookies, null, 2));
  };

  readCookies = async (filepath) => {
    let cookies = [];

    try {
      const cookiesString = await fs.promises.readFile(filepath);
      cookies = JSON.parse(cookiesString);
    } catch (err) {
      await fs.promises.writeFile(filepath, JSON.stringify([], null, 2));
      cookies = [];
    }

    await this.context.addCookies(cookies);
  };

  abortImageRequests = async () => {
    await this.context.route('**', (route) => {
      if (route.request().resourceType() === 'image') {
        route.abort();
      } else {
        route.continue();
      }
    });
  };

  // delay for n seconds
  customDelay = async (n) => {
    await this.page.waitForTimeout(n * 1000);
  };

  // delay between 1-3 seconds
  fastDelay = async () => {
    const delay = (Math.floor(Math.random() * 2) + 1) * 1000;
    await this.page.waitForTimeout(delay);
  };

  doScreenshot = async (filepath) => {
    await this.page.screenshot({
      path: filepath,
      fullPage: true
    });
  };

  gotoWait = async (url, timeout = 5) => {
    await this.page.goto(url, {
      waitUntil: 'load',
      timeout: timeout * 1000
    });
  };

  /**
   * @param {string} text
   * @param {number} timeout
   * @return {Promise<Locator>}
   */
  getButtonWithText = async (text, timeout = 3) => {
    const button = this.page.locator(`button:text("${text}")`);
    try {
      await button.waitFor({
        state: 'visible',
        timeout: timeout * 1000
      });
      return button;
    } catch {
      return undefined;
    }
  };

  waitForLocatorElement = async (locator, timeout = 5) => {
    try {
      const button = this.page.locator(locator);
      await button.waitFor({
        state: 'visible',
        timeout: timeout * 1000
      });
      return button;
    } catch {
      return undefined;
    }
  };

  waitForSelectorElement = async (selector, timeout = 5) => {
    try {
      const elem = await this.page.waitForSelector(selector, {
        timeout: timeout * 1000
      });
      return elem;
    } catch {
      return undefined;
    }
  };

  waitForFileDownload = async (timeout = 120) => {
    try {
      const download = await this.page.waitForEvent('download', {
        timeout: timeout * 1000
      });
      const downloadPath = await download.path();
      const filename = download.suggestedFilename();
      const filePath = path.dirname(downloadPath) + '/' + filename;
      await sleep(5000);
      fs.renameSync(downloadPath, filePath);
      await sleep(2000);
      return [filePath, true];
    } catch (err) {
      console.log(err);
      return ['', false];
    }
  };
}

module.exports = Scrapper;
