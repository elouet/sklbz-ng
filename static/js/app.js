// Example JavaScript file
console.log('Welcome to sklbz-ng!');

// Function to fetch articles from API
async function fetchArticles() {
    try {
        const response = await fetch('/api/articles');
        const articles = await response.json();
        console.log('Articles:', articles);
        return articles;
    } catch (error) {
        console.error('Error fetching articles:', error);
        return [];
    }
}

// Function to create a new article
async function createArticle(articleData) {
    try {
        const response = await fetch('/api/articles', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(articleData),
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const newArticle = await response.json();
        console.log('Created article:', newArticle);
        return newArticle;
    } catch (error) {
        console.error('Error creating article:', error);
        throw error;
    }
}

// Example usage
// fetchArticles();
// createArticle({
//     title: 'My First Article',
//     content: 'This is the content of my first article.',
//     author: 'John Doe',
//     visibility: 'public',
//     categories: ['technology', 'web'],
//     tags: ['javascript', 'api', 'blog']
// });
