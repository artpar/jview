#import <Cocoa/Cocoa.h>
#include "list.h"
#include "stackview.h"

void* JVCreateList(const char* justify, const char* align, int gap, int padding) {
    // List defaults to stretch alignment so children fill width
    const char* effectiveAlign = (align && strlen(align) > 0) ? align : "stretch";
    const char* effectiveJustify = (justify && strlen(justify) > 0) ? justify : "start";
    return JVCreateStackView(false, effectiveJustify, effectiveAlign, gap, padding);
}

void JVUpdateList(void* handle, const char* justify, const char* align, int gap, int padding) {
    const char* effectiveAlign = (align && strlen(align) > 0) ? align : "stretch";
    const char* effectiveJustify = (justify && strlen(justify) > 0) ? justify : "start";
    JVUpdateStackView(handle, effectiveJustify, effectiveAlign, gap, padding);
}

void JVListSetChildren(void* handle, void** children, int count) {
    JVStackViewSetChildren(handle, children, count);
}
