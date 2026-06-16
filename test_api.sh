#!/bin/bash

# Test script for sklbz-ng API
# This script assumes the server is running on localhost:8080

SERVER_URL="http://localhost:8080"

echo "Testing sklbz-ng API..."
echo "Server URL: $SERVER_URL"
echo ""

# Test 1: Health check
echo "1. Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s $SERVER_URL/health)
echo "Response: $HEALTH_RESPONSE"
if [[ $HEALTH_RESPONSE == *"healthy"* ]]; then
    echo "✓ Health check passed"
else
    echo "✗ Health check failed"
fi
echo ""

# Test 2: Root endpoint
echo "2. Testing root endpoint..."
ROOT_RESPONSE=$(curl -s $SERVER_URL/)
echo "Response: $ROOT_RESPONSE"
if [[ $ROOT_RESPONSE == *"Welcome to sklbz-ng API"* ]]; then
    echo "✓ Root endpoint passed"
else
    echo "✗ Root endpoint failed"
fi
echo ""

# Test 3: Create an article
echo "3. Creating a test article..."
CREATE_RESPONSE=$(curl -s -X POST $SERVER_URL/api/articles \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Article",
    "content": "This is a test article content.",
    "author": "Test Author",
    "visibility": "public",
    "categories": ["test", "api"],
    "tags": ["test", "curl"]
  }')
echo "Response: $CREATE_RESPONSE"
ARTICLE_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
if [[ -n $ARTICLE_ID ]]; then
    echo "✓ Article created with ID: $ARTICLE_ID"
else
    echo "✗ Failed to create article"
fi
echo ""

# Test 4: List articles
echo "4. Listing articles..."
LIST_RESPONSE=$(curl -s $SERVER_URL/api/articles)
echo "Response: $LIST_RESPONSE"
if [[ $LIST_RESPONSE == *"Test Article"* ]]; then
    echo "✓ Article found in list"
else
    echo "✗ Article not found in list"
fi
echo ""

# Test 5: Get specific article (if we got an ID)
if [[ -n $ARTICLE_ID ]]; then
    echo "5. Getting specific article..."
    GET_RESPONSE=$(curl -s $SERVER_URL/api/articles/$ARTICLE_ID)
    echo "Response: $GET_RESPONSE"
    if [[ $GET_RESPONSE == *"Test Article"* ]]; then
        echo "✓ Article retrieved successfully"
    else
        echo "✗ Failed to retrieve article"
    fi
    echo ""

    # Test 6: Update article
    echo "6. Updating article..."
    UPDATE_RESPONSE=$(curl -s -X PUT $SERVER_URL/api/articles/$ARTICLE_ID \
      -H "Content-Type: application/json" \
      -d '{
        "title": "Updated Test Article",
        "content": "This is updated content."
      }')
    echo "Response: $UPDATE_RESPONSE"
    if [[ $UPDATE_RESPONSE == *"Updated Test Article"* ]]; then
        echo "✓ Article updated successfully"
    else
        echo "✗ Failed to update article"
    fi
    echo ""

    # Test 7: Delete article
    echo "7. Deleting article..."
    DELETE_RESPONSE=$(curl -s -X DELETE $SERVER_URL/api/articles/$ARTICLE_ID)
    echo "Response status: $DELETE_RESPONSE"
    if [[ -z $DELETE_RESPONSE ]]; then
        echo "✓ Article deleted successfully"
    else
        echo "✗ Failed to delete article"
    fi
    echo ""

    # Test 8: Verify deletion
    echo "8. Verifying deletion..."
    VERIFY_RESPONSE=$(curl -s $SERVER_URL/api/articles/$ARTICLE_ID)
    if [[ $VERIFY_RESPONSE == *"not found"* ]]; then
        echo "✓ Article successfully deleted"
    else
        echo "✗ Article still exists"
    fi
fi

# Test 9: Static files
echo "9. Testing static files..."
JS_RESPONSE=$(curl -s -I $SERVER_URL/static/js/app.js | head -1)
CSS_RESPONSE=$(curl -s -I $SERVER_URL/static/css/style.css | head -1)

if [[ $JS_RESPONSE == *"200 OK"* ]]; then
    echo "✓ JavaScript file accessible"
else
    echo "✗ JavaScript file not accessible"
fi

if [[ $CSS_RESPONSE == *"200 OK"* ]]; then
    echo "✓ CSS file accessible"
else
    echo "✗ CSS file not accessible"
fi

echo ""
echo "API testing complete!"
