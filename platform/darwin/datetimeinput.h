#ifndef JVIEW_DATETIMEINPUT_H
#define JVIEW_DATETIMEINPUT_H

#include <stdint.h>
#include <stdbool.h>

void* JVCreateDateTimeInput(bool enableDate, bool enableTime, const char* value, uint64_t callbackID);
void JVUpdateDateTimeInput(void* handle, bool enableDate, bool enableTime, const char* value);

#endif
