const fs = require('fs');
const path = require('path');
const OSS = require('ali-oss');
const { promisify } = require('util');
const readdir = promisify(fs.readdir);
const readFile = promisify(fs.readFile);

const client = new OSS({
    region: '', // Replace with your region
    accessKeyId: '', // Replace with your access key id
    accessKeySecret: '', // Replace with your access key secret
    bucket: 'user-avatar-test' // Replace with your bucket name
});

// Read and parse the JSON file
const data = JSON.parse(fs.readFileSync(path.join(__dirname, 'data.json'), 'utf8'));

// Function to get all image files from a directory
async function getImageFiles(directory) {
    try {
        const files = await readdir(path.join(__dirname, directory));
        return files.filter(file => /\.(jpg|jpeg|png|gif)$/i.test(file));
    } catch (error) {
        console.error(`Error reading directory ${directory}:`, error);
        return [];
    }
}

// Function to upload a file to Aliyun OSS
async function uploadToOSS(filePath, fileName) {
    try {
        const result = await client.put(fileName, filePath);
        return result.url;
    } catch (error) {
        console.error(`Error uploading ${fileName}:`, error);
        return null;
    }
}

// Function to get random user ID from database
async function getRandomUserId() {
    // In a real application, you would query the database here
    // For this example, we'll use a hardcoded list of user IDs
    const userIds = [1, 2, 3, 4, 5, 6]; // Replace with actual user IDs from your database
    return userIds[Math.floor(Math.random() * userIds.length)];
}

// Main function to process the data
async function processData() {
    const sqlStatements = [];
    
    for (const guide of data) {
        // Get all images from the specified directory
        const imageFiles = await getImageFiles(guide.image_path);
        const imageUrls = [];
        
        // Upload each image to Aliyun OSS
        for (const file of imageFiles) {
            const filePath = path.join(__dirname, guide.image_path, file);
            const fileName = `guides/${Date.now()}-${file}`;
            const url = await uploadToOSS(filePath, fileName);
            if (url) {
                imageUrls.push(url);
            }
        }
        
        // Get random user ID
        const userId = await getRandomUserId();
        
        // Create SQL insert statement
        const sql = `INSERT INTO travel_guides (title, content, images, user_id, published_at) VALUES (
            '${guide.title.replace(/'/g, "''")}',
            '${guide.desc.replace(/'/g, "''")}',
            '${JSON.stringify(imageUrls)}',
            ${userId},
            NOW()
        );`;
        
        sqlStatements.push(sql);
    }
    
    // Write SQL statements to file
    fs.writeFileSync(
        path.join(__dirname, 'insert_guides.sql'),
        sqlStatements.join('\n\n'),
        'utf8'
    );
    
    console.log('Processing completed. SQL statements written to insert_guides.sql');
}

// Execute the main function
processData().catch(console.error);
