#ifndef JVIEW_CARD_H
#define JVIEW_CARD_H

void* JVCreateCard(const char* title, const char* subtitle, int padding);
void JVUpdateCard(void* handle, const char* title, const char* subtitle, int padding);
void JVCardSetChildren(void* handle, void** children, int count);

#endif
