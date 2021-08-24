import puppeteer from "puppeteer";

describe("App.js", () => {
  let browser : any;
  let page : any;

  beforeAll(async () => {
    browser = await puppeteer.launch();
    page = await browser.newPage();
  });

  it("contains the welcome text", async () => {
    await page.goto("http://localhost:5000/signup");
    await page.waitForSelector("#login-form");
    // const text = await page.$eval(".App-welcome-text", (e) => e.textContent);
    // expect(text).toContain("Edit src/App.js and save to reload.");
  });

  afterAll(() => browser.close());
});