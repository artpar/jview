#ifndef JVIEW_LIST_H
#define JVIEW_LIST_H

void* JVCreateList(const char* justify, const char* align, int gap, int padding);
void JVUpdateList(void* handle, const char* justify, const char* align, int gap, int padding);
void JVListSetChildren(void* handle, void** children, int count);

#endif
