#ifndef JVIEW_CHECKBOX_H
#define JVIEW_CHECKBOX_H

#include <stdint.h>
#include <stdbool.h>

void* JVCreateCheckBox(const char* label, bool checked, uint64_t callbackID);
void JVUpdateCheckBox(void* handle, const char* label, bool checked);

#endif
