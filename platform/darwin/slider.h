#ifndef JVIEW_SLIDER_H
#define JVIEW_SLIDER_H

#include <stdint.h>

void* JVCreateSlider(double min, double max, double step, double value, uint64_t callbackID);
void JVUpdateSlider(void* handle, double min, double max, double step, double value);

#endif
