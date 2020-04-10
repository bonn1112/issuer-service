const fs = require('fs');
const puppeteer = require('puppeteer');

const args = process.argv;
if (args.length < 4) {
    console.error("filepath's is required");
    return
}

const htmlFilepath = args[2];
const pdfFilepath = args[3];

(async () => {
    const buf = fs.readFileSync(htmlFilepath, {encoding: 'utf-8'});
    const html = buf.toString('ascii');

    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    await page.setContent(html, {waitUntil: 'networkidle2'});
    await page.pdf({path: pdfFilepath, format: 'A4'});
    await browser.close()
})();
