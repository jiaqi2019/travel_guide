const mysql = require('mysql2/promise');
const fs = require('fs').promises;
const path = require('path');

async function convertGuides() {
    // 读取data.json文件
    const data = JSON.parse(await fs.readFile('data.json', 'utf8'));

    // 数据库连接配置
    const dbConfig = {
        host: 'localhost',
        user: 'root',
        password: '123456',
        database: 'travel_guide',
        charset: 'utf8mb4'
    };

    // 连接数据库
    const connection = await mysql.createConnection(dbConfig);

    try {
        // 获取所有用户ID
        const [users] = await connection.execute('SELECT id FROM users');
        const userIds = users.map(user => user.id);

        // 获取所有标签ID
        const [tags] = await connection.execute('SELECT id FROM tags');
        const tagIds = tags.map(tag => tag.id);

        // 准备SQL语句
        const insertGuideSql = `
            INSERT INTO travel_guides (user_id, title, content, images, published_at, created_at, updated_at)
            VALUES (?, ?, ?, ?, FROM_UNIXTIME(?), FROM_UNIXTIME(?), FROM_UNIXTIME(?))
        `;

        const insertGuideTagSql = `
            INSERT INTO guide_tags (guide_id, tag_id)
            VALUES (?, ?)
        `;

        // 处理数据
        let totalGuides = 0;
        for (const [tag, items] of Object.entries(data)) {
            console.log(`Processing tag: ${tag}`);
            
            for (const item of items) {
                try {
                    const noteCard = item.note_card;
                    if (!noteCard) {
                        console.log('Skipping item without note_card');
                        continue;
                    }

                    // 提取图片URL
                    const images = noteCard.image_list.map(img => img.url_default);
                    const imagesJson = JSON.stringify(images);

                    // 随机选择一个用户ID
                    const userId = userIds[Math.floor(Math.random() * userIds.length)];

                    // 当前时间戳（秒）
                    const now = Math.floor(Date.now() / 1000);

                    // 插入游记数据
                    const [result] = await connection.execute(insertGuideSql, [
                        userId,
                        noteCard.title,
                        noteCard.desc,
                        imagesJson,
                        now,
                        now,
                        now
                    ]);

                    // 获取刚插入的游记ID
                    const guideId = result.insertId;

                    // 查找对应的标签ID
                    const [tagRows] = await connection.execute(
                        'SELECT id FROM tags WHERE name = ?',
                        [tag]
                    );

                    if (tagRows.length > 0) {
                        // 插入标签关联
                        await connection.execute(insertGuideTagSql, [
                            guideId,
                            tagRows[0].id
                        ]);
                    }

                    totalGuides++;
                    console.log(`Processed guide: ${noteCard.title}`);
                } catch (error) {
                    console.error('Error processing item:', error);
                    continue;
                }
            }
        }

        console.log(`数据转换完成！共处理 ${totalGuides} 条游记`);
    } catch (error) {
        console.error('转换过程中发生错误：', error);
    } finally {
        // 关闭连接
        await connection.end();
    }
}

// 执行转换
convertGuides(); 