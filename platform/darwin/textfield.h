#ifndef JVIEW_TEXTFIELD_H
#define JVIEW_TEXTFIELD_H

#include <stdint.h>
#include <stdbool.h>

void* JVCreateTextField(const char* placeholder, const char* value,
                         const char* inputType, bool readOnly, uint64_t callbackID);
void JVUpdateTextField(void* handle, const char* placeholder, const char* value,
                        const char* inputType, bool readOnly);
void JVSetTextFieldErrors(void* handle, const char** errors, int count);

#endif
