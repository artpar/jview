#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <ctype.h>

// Proper native C functions with real signatures — no JSON wrappers needed.
// libffi calls these directly with the correct C types.

double math_add(double a, double b) {
    return a + b;
}

int string_length(const char *s) {
    return (int)strlen(s);
}

// Returns a pointer to a static buffer (caller must NOT free).
const char* string_reverse(const char *s) {
    static char buf[4096];
    size_t len = strlen(s);
    if (len >= sizeof(buf)) len = sizeof(buf) - 1;
    for (size_t i = 0; i < len; i++) {
        buf[i] = s[len - 1 - i];
    }
    buf[len] = '\0';
    return buf;
}

// Returns a pointer to a static buffer (caller must NOT free).
const char* string_upper(const char *s) {
    static char buf[4096];
    size_t len = strlen(s);
    if (len >= sizeof(buf)) len = sizeof(buf) - 1;
    for (size_t i = 0; i < len; i++) {
        buf[i] = (char)toupper((unsigned char)s[i]);
    }
    buf[len] = '\0';
    return buf;
}

// echo: returns the input string as-is (identity function for testing).
const char* echo(const char *s) {
    return s;
}

void* alloc_buffer(int size) {
    return malloc((size_t)size);
}

void free_buffer(void *ptr) {
    free(ptr);
}

int int_add(int a, int b) {
    return a + b;
}

float float_add(float a, float b) {
    return a + b;
}
