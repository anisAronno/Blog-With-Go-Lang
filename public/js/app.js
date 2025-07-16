// public/js/app.js - Main JavaScript functionality for Go Blog

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
});

// Main app initialization
function initializeApp() {
    initializeToasts();
    initializeFormValidation();
    initializeAnimations();
    initializeTheme();
    initializeSearchFeatures();
    console.log('🚀 Go Blog application initialized');
}

// Toast notification system
function initializeToasts() {
    // Auto-hide alerts with data-auto-hide attribute
    const alerts = document.querySelectorAll('[data-auto-hide]');
    alerts.forEach(alert => {
        setTimeout(() => {
            fadeOut(alert);
        }, 5000);
    });
}

// Toast notification function
function showToast(message, type = 'info', duration = 3000) {
    const toast = document.createElement('div');
    toast.className = `fixed top-4 right-4 z-50 px-4 py-3 rounded-lg shadow-lg text-white mb-2 transform transition-all duration-300 ${
        type === 'success' ? 'bg-green-500' : 
        type === 'error' ? 'bg-red-500' : 
        type === 'warning' ? 'bg-yellow-500' : 'bg-blue-500'
    }`;
    
    toast.innerHTML = `
        <div class="flex items-center">
            <i class="fas fa-${getToastIcon(type)} mr-2"></i>
            <span>${message}</span>
            <button onclick="this.parentElement.parentElement.remove()" class="ml-4 text-white hover:text-gray-200">
                <i class="fas fa-times"></i>
            </button>
        </div>
    `;
    
    document.body.appendChild(toast);
    
    // Animate in
    setTimeout(() => {
        toast.style.transform = 'translateX(0)';
    }, 100);
    
    // Auto remove
    setTimeout(() => {
        if (toast.parentElement) {
            toast.style.transform = 'translateX(100%)';
            setTimeout(() => {
                if (toast.parentElement) {
                    toast.remove();
                }
            }, 300);
        }
    }, duration);
}

// Get appropriate icon for toast type
function getToastIcon(type) {
    switch(type) {
        case 'success': return 'check-circle';
        case 'error': return 'exclamation-circle';
        case 'warning': return 'exclamation-triangle';
        default: return 'info-circle';
    }
}

// Form validation
function initializeFormValidation() {
    // Real-time email validation
    const emailInputs = document.querySelectorAll('input[type="email"]');
    emailInputs.forEach(input => {
        input.addEventListener('blur', validateEmail);
        input.addEventListener('input', clearValidationError);
    });
    
    // Real-time password validation
    const passwordInputs = document.querySelectorAll('input[type="password"]');
    passwordInputs.forEach(input => {
        input.addEventListener('input', validatePassword);
    });
    
    // Form submission validation
    const forms = document.querySelectorAll('form');
    forms.forEach(form => {
        form.addEventListener('submit', handleFormSubmission);
    });
}

// Email validation
function validateEmail(event) {
    const input = event.target;
    const email = input.value.trim();
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    
    if (email && !emailRegex.test(email)) {
        showValidationError(input, 'Please enter a valid email address');
        return false;
    }
    
    clearValidationError(input);
    return true;
}

// Password validation
function validatePassword(event) {
    const input = event.target;
    const password = input.value;
    const minLength = 6;
    
    if (password.length > 0 && password.length < minLength) {
        showValidationError(input, `Password must be at least ${minLength} characters long`);
        return false;
    }
    
    clearValidationError(input);
    
    // Check for password confirmation match
    const confirmInput = document.querySelector('input[name="password_confirmation"]');
    if (confirmInput && confirmInput.value && input.name === 'password') {
        validatePasswordConfirmation(confirmInput);
    }
    
    return true;
}

// Password confirmation validation
function validatePasswordConfirmation(confirmInput) {
    const passwordInput = document.querySelector('input[name="password"]');
    if (!passwordInput) return true;
    
    if (confirmInput.value !== passwordInput.value) {
        showValidationError(confirmInput, 'Passwords do not match');
        return false;
    }
    
    clearValidationError(confirmInput);
    return true;
}

// Show validation error
function showValidationError(input, message) {
    clearValidationError(input);
    
    input.classList.add('border-red-500');
    const errorDiv = document.createElement('div');
    errorDiv.className = 'validation-error text-red-500 text-sm mt-1';
    errorDiv.textContent = message;
    
    input.parentNode.appendChild(errorDiv);
}

// Clear validation error
function clearValidationError(input) {
    if (typeof input === 'object' && input.target) {
        input = input.target;
    }
    
    input.classList.remove('border-red-500');
    const existingError = input.parentNode.querySelector('.validation-error');
    if (existingError) {
        existingError.remove();
    }
}

// Handle form submission
function handleFormSubmission(event) {
    const form = event.target;
    const submitButton = form.querySelector('button[type="submit"]');
    
    // Show loading state
    if (submitButton) {
        const originalText = submitButton.innerHTML;
        submitButton.innerHTML = '<span class="spinner mr-2"></span>Processing...';
        submitButton.disabled = true;
        
        // Reset button after a delay if form submission fails
        setTimeout(() => {
            if (submitButton.disabled) {
                submitButton.innerHTML = originalText;
                submitButton.disabled = false;
            }
        }, 10000);
    }
    
    // Validate all form fields
    let isValid = true;
    const emailInputs = form.querySelectorAll('input[type="email"]');
    const passwordInputs = form.querySelectorAll('input[type="password"]');
    const requiredInputs = form.querySelectorAll('input[required], textarea[required]');
    
    // Validate required fields
    requiredInputs.forEach(input => {
        if (!input.value.trim()) {
            showValidationError(input, 'This field is required');
            isValid = false;
        }
    });
    
    // Validate email fields
    emailInputs.forEach(input => {
        if (!validateEmail({target: input})) {
            isValid = false;
        }
    });
    
    // Validate password fields
    passwordInputs.forEach(input => {
        if (!validatePassword({target: input})) {
            isValid = false;
        }
    });
    
    // Check password confirmation
    const confirmInput = form.querySelector('input[name="password_confirmation"]');
    if (confirmInput && !validatePasswordConfirmation(confirmInput)) {
        isValid = false;
    }
    
    if (!isValid) {
        event.preventDefault();
        if (submitButton) {
            submitButton.innerHTML = originalText;
            submitButton.disabled = false;
        }
        showToast('Please fix the errors below', 'error');
    }
}

// Animation utilities
function initializeAnimations() {
    // Add fade-in animation to cards
    const cards = document.querySelectorAll('.shadow-card, .dashboard-card');
    cards.forEach((card, index) => {
        card.style.opacity = '0';
        card.style.transform = 'translateY(20px)';
        
        setTimeout(() => {
            card.style.transition = 'all 0.5s ease';
            card.style.opacity = '1';
            card.style.transform = 'translateY(0)';
        }, index * 100);
    });
    
    // Add hover effects
    const hoverElements = document.querySelectorAll('.card-hover');
    hoverElements.forEach(element => {
        element.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-5px)';
        });
        
        element.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0)';
        });
    });
}

// Fade out animation
function fadeOut(element) {
    element.style.transition = 'opacity 0.3s ease';
    element.style.opacity = '0';
    setTimeout(() => {
        if (element.parentElement) {
            element.remove();
        }
    }, 300);
}

// Theme management
function initializeTheme() {
    // Check for saved theme preference or default to light mode
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
    
    // Theme toggle functionality (if theme toggle button exists)
    const themeToggle = document.getElementById('theme-toggle');
    if (themeToggle) {
        themeToggle.addEventListener('click', toggleTheme);
    }
}

// Toggle theme
function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    
    showToast(`Switched to ${newTheme} mode`, 'info');
}

// Search functionality
function initializeSearchFeatures() {
    const searchInput = document.getElementById('search-input');
    if (searchInput) {
        searchInput.addEventListener('input', debounce(handleSearch, 300));
    }
}

// Search handler
function handleSearch(event) {
    const query = event.target.value.trim();
    if (query.length < 2) {
        clearSearchResults();
        return;
    }
    
    // In a real application, this would make an AJAX call to search endpoint
    console.log('Searching for:', query);
    
    // Mock search results for demonstration
    showSearchResults([
        { title: 'Sample Blog Post', url: '/blog/1' },
        { title: 'Another Post', url: '/blog/2' }
    ]);
}

// Show search results
function showSearchResults(results) {
    let resultsContainer = document.getElementById('search-results');
    if (!resultsContainer) {
        resultsContainer = document.createElement('div');
        resultsContainer.id = 'search-results';
        resultsContainer.className = 'absolute top-full left-0 right-0 bg-white border border-gray-200 rounded-md shadow-lg z-50 mt-1';
        document.getElementById('search-input').parentNode.appendChild(resultsContainer);
    }
    
    if (results.length === 0) {
        resultsContainer.innerHTML = '<div class="p-4 text-gray-500">No results found</div>';
        return;
    }
    
    const resultsHTML = results.map(result => `
        <a href="${result.url}" class="block p-3 hover:bg-gray-50 border-b border-gray-100 last:border-b-0">
            <div class="font-medium text-gray-900">${result.title}</div>
        </a>
    `).join('');
    
    resultsContainer.innerHTML = resultsHTML;
}

// Clear search results
function clearSearchResults() {
    const resultsContainer = document.getElementById('search-results');
    if (resultsContainer) {
        resultsContainer.remove();
    }
}

// Utility: Debounce function
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Confirm delete function
function confirmDelete(message = 'Are you sure you want to delete this item?') {
    return confirm(message);
}

// Copy to clipboard function
function copyToClipboard(text) {
    if (navigator.clipboard) {
        navigator.clipboard.writeText(text).then(() => {
            showToast('Copied to clipboard!', 'success');
        }).catch(() => {
            showToast('Failed to copy to clipboard', 'error');
        });
    } else {
        // Fallback for older browsers
        const textArea = document.createElement('textarea');
        textArea.value = text;
        document.body.appendChild(textArea);
        textArea.select();
        try {
            document.execCommand('copy');
            showToast('Copied to clipboard!', 'success');
        } catch (err) {
            showToast('Failed to copy to clipboard', 'error');
        }
        document.body.removeChild(textArea);
    }
}

// Auto-resize textarea
function autoResizeTextarea(textarea) {
    textarea.style.height = 'auto';
    textarea.style.height = textarea.scrollHeight + 'px';
}

// Initialize auto-resize for textareas
document.addEventListener('DOMContentLoaded', function() {
    const textareas = document.querySelectorAll('textarea');
    textareas.forEach(textarea => {
        textarea.addEventListener('input', function() {
            autoResizeTextarea(this);
        });
        // Initial resize
        autoResizeTextarea(textarea);
    });
});

// Keyboard shortcuts
document.addEventListener('keydown', function(event) {
    // Ctrl/Cmd + K for search
    if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
        event.preventDefault();
        const searchInput = document.getElementById('search-input');
        if (searchInput) {
            searchInput.focus();
        }
    }
    
    // Escape to close modals/dropdowns
    if (event.key === 'Escape') {
        clearSearchResults();
        // Close any open modals or dropdowns
        const openDropdowns = document.querySelectorAll('.dropdown-open');
        openDropdowns.forEach(dropdown => {
            dropdown.classList.remove('dropdown-open');
        });
    }
});

// Export functions for global use
window.showToast = showToast;
window.confirmDelete = confirmDelete;
window.copyToClipboard = copyToClipboard;
