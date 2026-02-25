#ifndef JVIEW_SEARCHFIELD_H
#define JVIEW_SEARCHFIELD_H

#include <stdint.h>

void* JVCreateSearchField(const char* placeholder, const char* value, uint64_t callbackID);
void JVUpdateSearchField(void* handle, const char* placeholder, const char* value);

#endif
