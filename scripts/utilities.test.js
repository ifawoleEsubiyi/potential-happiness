/**
 * Test file for utilities.js
 * 
 * Run these tests in a browser console or with a test runner like Jest
 */

// Mock document object for Node.js testing
if (typeof document === 'undefined') {
    global.document = {
        querySelectorAll: function() {
            return [];
        }
    };
}

// Import functions if running in Node.js
let calculateDaysBetweenDates, highlightImagesWithoutAlt, removeImageHighlighting;
if (typeof require !== 'undefined') {
    const utils = require('./utilities.js');
    calculateDaysBetweenDates = utils.calculateDaysBetweenDates;
    highlightImagesWithoutAlt = utils.highlightImagesWithoutAlt;
    removeImageHighlighting = utils.removeImageHighlighting;
}

/**
 * Test suite for calculateDaysBetweenDates
 */
function testCalculateDaysBetweenDates() {
    console.log('Testing calculateDaysBetweenDates...');
    
    const tests = [
        {
            name: 'Same date should return 0',
            begin: '2024-01-01',
            end: '2024-01-01',
            expected: 0
        },
        {
            name: 'One day difference',
            begin: '2024-01-01',
            end: '2024-01-02',
            expected: 1
        },
        {
            name: 'One week difference',
            begin: '2024-01-01',
            end: '2024-01-08',
            expected: 7
        },
        {
            name: 'End before begin (negative result)',
            begin: '2024-01-10',
            end: '2024-01-01',
            expected: -9
        },
        {
            name: 'Leap year - across Feb 29',
            begin: '2024-02-28',
            end: '2024-03-01',
            expected: 2 // 2024 is a leap year
        },
        {
            name: 'One year difference',
            begin: '2024-01-01',
            end: '2025-01-01',
            expected: 366 // 2024 is a leap year
        },
        {
            name: 'Date objects as input',
            begin: new Date('2024-01-01'),
            end: new Date('2024-01-10'),
            expected: 9
        }
    ];
    
    let passed = 0;
    let failed = 0;
    
    tests.forEach(test => {
        try {
            const result = calculateDaysBetweenDates(test.begin, test.end);
            if (result === test.expected) {
                console.log(`✓ ${test.name}: PASSED`);
                passed++;
            } else {
                console.error(`✗ ${test.name}: FAILED - Expected ${test.expected}, got ${result}`);
                failed++;
            }
        } catch (error) {
            console.error(`✗ ${test.name}: ERROR - ${error.message}`);
            failed++;
        }
    });
    
    // Test error handling
    try {
        calculateDaysBetweenDates('invalid-date', '2024-01-01');
        console.error('✗ Invalid date should throw error: FAILED');
        failed++;
    } catch (error) {
        console.log('✓ Invalid date throws error: PASSED');
        passed++;
    }
    
    console.log(`\nResults: ${passed} passed, ${failed} failed\n`);
    return { passed, failed };
}

/**
 * Test suite for highlightImagesWithoutAlt
 * Note: This requires a DOM environment to test properly
 */
function testHighlightImagesWithoutAlt() {
    console.log('Testing highlightImagesWithoutAlt...');
    
    if (typeof document === 'undefined' || typeof document.createElement === 'undefined') {
        console.log('Skipping DOM tests - no browser environment detected');
        return { passed: 0, failed: 0, skipped: true };
    }
    
    // Create test images
    const testContainer = document.createElement('div');
    testContainer.id = 'test-container';
    
    // Image with alt text (should not be highlighted)
    const imgWithAlt = document.createElement('img');
    imgWithAlt.src = 'test1.jpg';
    imgWithAlt.alt = 'Test image with alt text';
    testContainer.appendChild(imgWithAlt);
    
    // Image without alt attribute (should be highlighted)
    const imgWithoutAlt = document.createElement('img');
    imgWithoutAlt.src = 'test2.jpg';
    testContainer.appendChild(imgWithoutAlt);
    
    // Image with empty alt (should be highlighted)
    const imgWithEmptyAlt = document.createElement('img');
    imgWithEmptyAlt.src = 'test3.jpg';
    imgWithEmptyAlt.alt = '';
    testContainer.appendChild(imgWithEmptyAlt);
    
    // Image with whitespace-only alt (should be highlighted)
    const imgWithWhitespaceAlt = document.createElement('img');
    imgWithWhitespaceAlt.src = 'test4.jpg';
    imgWithWhitespaceAlt.alt = '   ';
    testContainer.appendChild(imgWithWhitespaceAlt);
    
    document.body.appendChild(testContainer);
    
    // Run the function
    const highlightedImages = highlightImagesWithoutAlt();
    
    let passed = 0;
    let failed = 0;
    
    // Test: Should highlight 3 images (without alt, empty alt, whitespace alt)
    if (highlightedImages.length === 3) {
        console.log('✓ Correct number of images highlighted: PASSED');
        passed++;
    } else {
        console.error(`✗ Expected 3 images highlighted, got ${highlightedImages.length}: FAILED`);
        failed++;
    }
    
    // Test: Image with alt should not have border
    if (!imgWithAlt.style.border || imgWithAlt.style.border === '') {
        console.log('✓ Image with alt text not highlighted: PASSED');
        passed++;
    } else {
        console.error('✗ Image with alt text should not be highlighted: FAILED');
        failed++;
    }
    
    // Test: Image without alt should have red border
    if (imgWithoutAlt.style.border === '3px solid red') {
        console.log('✓ Image without alt has red border: PASSED');
        passed++;
    } else {
        console.error(`✗ Expected red border, got: ${imgWithoutAlt.style.border}: FAILED`);
        failed++;
    }
    
    // Clean up
    document.body.removeChild(testContainer);
    
    console.log(`\nResults: ${passed} passed, ${failed} failed\n`);
    return { passed, failed };
}

/**
 * Run all tests
 */
function runAllTests() {
    console.log('=== Running All Tests ===\n');
    
    const results = {
        calculateDays: testCalculateDaysBetweenDates(),
        highlightImages: testHighlightImagesWithoutAlt()
    };
    
    const totalPassed = results.calculateDays.passed + results.highlightImages.passed;
    const totalFailed = results.calculateDays.failed + results.highlightImages.failed;
    
    console.log('=== Overall Results ===');
    console.log(`Total: ${totalPassed} passed, ${totalFailed} failed`);
    
    return totalFailed === 0;
}

// Export for Node.js
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        testCalculateDaysBetweenDates,
        testHighlightImagesWithoutAlt,
        runAllTests
    };
}

// Auto-run tests if executed directly in Node.js
if (typeof require !== 'undefined' && require.main === module) {
    runAllTests();
}
