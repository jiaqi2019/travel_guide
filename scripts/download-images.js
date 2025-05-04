const fs = require('fs');
const path = require('path');
const axios = require('axios');
const cheerio = require('cheerio');

// Read the data.json file
const dataPath = path.join(__dirname, '../download-img/data.json');
const outputDir = path.join(__dirname, '../download-img/images');

// Create output directory if it doesn't exist
if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
}

async function downloadImages() {
    try {
        // Read the JSON file
        const data = JSON.parse(fs.readFileSync(dataPath, 'utf8'));
        
        for (const url of data) {
            console.log(`Processing URL: ${url}`);
            
            // Fetch the page content
            const response = await axios.get(url);
            const $ = cheerio.load(response.data);
            
            // Find all image elements
            const images = $('img');
            
            // Download each image
            for (let i = 0; i < images.length; i++) {
                const imgUrl = $(images[i]).attr('src');
                if (imgUrl) {
                    try {
                        const imgResponse = await axios({
                            url: imgUrl,
                            responseType: 'arraybuffer'
                        });
                        
                        const fileName = `image_${Date.now()}_${i}.jpg`;
                        const filePath = path.join(outputDir, fileName);
                        
                        fs.writeFileSync(filePath, imgResponse.data);
                        console.log(`Downloaded: ${fileName}`);
                    } catch (error) {
                        console.error(`Error downloading image: ${imgUrl}`, error.message);
                    }
                }
            }
        }
    } catch (error) {
        console.error('Error:', error.message);
    }
}

downloadImages(); 