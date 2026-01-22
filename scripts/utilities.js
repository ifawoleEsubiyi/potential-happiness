/**
 * JavaScript Utility Functions
 * 
 * This file contains various utility functions for web and date operations.
 */

// Constants
const MILLISECONDS_PER_DAY = 1000 * 60 * 60 * 24;

/**
 * Calculate the number of days between two dates
 * 
 * @param {Date|string} begin - The start date
 * @param {Date|string} end - The end date
 * @returns {number} The number of days between the two dates
 * 
 * @example
 * calculateDaysBetweenDates('2024-01-01', '2024-01-10') // returns 9
 * calculateDaysBetweenDates(new Date('2024-01-01'), new Date('2024-01-10')) // returns 9
 */
function calculateDaysBetweenDates(begin, end) {
    // Convert inputs to Date objects if they aren't already
    const beginDate = begin instanceof Date ? begin : new Date(begin);
    const endDate = end instanceof Date ? end : new Date(end);
    
    // Validate dates
    if (isNaN(beginDate.getTime()) || isNaN(endDate.getTime())) {
        throw new Error('Invalid date provided');
    }
    
    // Calculate the difference in milliseconds
    const diffInMs = endDate.getTime() - beginDate.getTime();
    
    // Convert milliseconds to days
    const diffInDays = Math.floor(diffInMs / MILLISECONDS_PER_DAY);
    
    return diffInDays;
}

/**
 * Find all images without alternate text and give them a red border
 * 
 * This function searches the DOM for all <img> elements that don't have
 * an alt attribute or have an empty alt attribute, and applies a red border
 * to highlight them for accessibility review.
 * 
 * @param {string} [borderStyle='3px solid red'] - The CSS border style to apply
 * @returns {NodeList} The collection of images that were highlighted
 * 
 * @example
 * // Highlight all images without alt text with default red border
 * highlightImagesWithoutAlt();
 * 
 * // Highlight with custom border style
 * highlightImagesWithoutAlt('5px dashed orange');
 */
function highlightImagesWithoutAlt(borderStyle = '3px solid red') {
    // Get all image elements in the document
    const allImages = document.querySelectorAll('img');
    
    // Filter images without alt text or with empty alt text
    const imagesWithoutAlt = Array.from(allImages).filter(img => {
        // Check if alt attribute is missing or empty
        return !img.hasAttribute('alt') || img.getAttribute('alt').trim() === '';
    });
    
    // Apply red border to each image without alt text
    imagesWithoutAlt.forEach(img => {
        img.style.border = borderStyle;
    });
    
    // Return the collection for further processing if needed
    return imagesWithoutAlt;
}

/**
 * Remove highlighting from images (utility function to undo highlightImagesWithoutAlt)
 * 
 * @returns {void}
 */
function removeImageHighlighting() {
    const allImages = document.querySelectorAll('img');
    allImages.forEach(img => {
        img.style.border = '';
    });
}

// Export functions for use in Node.js/modules (if applicable)
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        calculateDaysBetweenDates,
        highlightImagesWithoutAlt,
        removeImageHighlighting
    };
}
