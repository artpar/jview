#ifndef JVIEW_SYSTEMEVENTS_H
#define JVIEW_SYSTEMEVENTS_H

// Start observing system appearance changes (light/dark mode).
// Calls GoSystemEvent("system.appearance", data) when it changes.
void JVStartAppearanceObserver(void);
void JVStopAppearanceObserver(void);

// Start observing power events (sleep/wake/battery).
// Calls GoSystemEvent("system.power.sleep", "{}") and GoSystemEvent("system.power.wake", "{}").
void JVStartPowerObserver(void);
void JVStopPowerObserver(void);

// Start observing display changes (connect/disconnect/resolution).
void JVStartDisplayObserver(void);
void JVStopDisplayObserver(void);

// Start observing locale/language changes.
void JVStartLocaleObserver(void);
void JVStopLocaleObserver(void);

// Start/stop polling clipboard for changes.
void JVStartClipboardObserver(int intervalMs);
void JVStopClipboardObserver(void);

// Start observing network reachability changes.
void JVStartNetworkObserver(void);
void JVStopNetworkObserver(void);

// Start observing accessibility changes (reduce motion, reduce transparency, etc.).
void JVStartAccessibilityObserver(void);
void JVStopAccessibilityObserver(void);

// Start observing thermal state changes.
void JVStartThermalObserver(void);
void JVStopThermalObserver(void);

#endif
