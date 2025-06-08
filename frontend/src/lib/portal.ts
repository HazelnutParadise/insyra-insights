export function portal(node: HTMLElement, target: HTMLElement | string = document.body) {
    let targetEl: HTMLElement | null;

    if (typeof target === 'string') {
        targetEl = document.querySelector(target);
        if (!targetEl) {
            // Fallback to document.body if target string doesn't match anything
            console.warn(`Portal target selector "${target}" not found. Falling back to document.body.`);
            targetEl = document.body;
        }
    } else {
        targetEl = target;
    }

    // Ensure targetEl is not null before appending (it could be if document.body was somehow null, though unlikely)
    if (targetEl) {
        targetEl.appendChild(node);
    } else {
        console.error("Portal target element (document.body or custom selector) is null. Node not appended.");
    }


    return {
        destroy() {
            // Ensure node and targetEl are still valid and node is a child of targetEl
            if (node && targetEl && node.parentNode === targetEl) {
                targetEl.removeChild(node);
            }
        }
    };
}
