type EventData = {
    [key: string]: unknown;
};

export function trackEvent(eventName: string, eventData: EventData) {
    if (typeof window === 'undefined' || !window.gtag) {
        console.error('gtag is not defined');
        return;
    }
    window.gtag('event', eventName, eventData);
}
