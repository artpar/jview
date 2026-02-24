#ifndef JVIEW_CHOICEPICKER_H
#define JVIEW_CHOICEPICKER_H

#include <stdint.h>

void* JVCreateChoicePicker(const char** labels, const char** values, int count,
                            const char* selected, uint64_t callbackID);
void JVUpdateChoicePicker(void* handle, const char** labels, const char** values, int count,
                           const char* selected);

#endif
