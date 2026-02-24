#ifndef JVIEW_STACKVIEW_H
#define JVIEW_STACKVIEW_H

#include <stdbool.h>

void* JVCreateStackView(bool horizontal, const char* justify, const char* align, int gap, int padding);
void JVUpdateStackView(void* handle, const char* justify, const char* align, int gap, int padding);
void JVStackViewSetChildren(void* handle, void** children, int count);

#endif
