let toastTimeout;
const toast = document.getElementById('toast');
const observer = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
        if (mutation.type === 'childList') {
            // Display toast
            toast.style.display = 'block';
            // Clear previous timeout if any
            if (toastTimeout) {
                clearTimeout(toastTimeout);
            }
            // Hide toast after 5 seconds
            toastTimeout = setTimeout(() => {
                toast.style.display = 'none';
            }, 5000);
        }
    }
});
observer.observe(toast, { childList: true });