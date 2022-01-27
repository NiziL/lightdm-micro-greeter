#define connect(name, func) g_signal_connect(greeter, name, G_CALLBACK(func), NULL)
#include "_cgo_export.h"

extern void authentication_complete_cb(LightDMGreeter *greeter);
extern void show_prompt_cb(LightDMGreeter *greeter, char *text, LightDMPromptType type);

void greeter_signal_connect(LightDMGreeter* greeter) {
    connect("authentication-complete", authentication_complete_cb);
    connect("show-prompt", show_prompt_cb);
}
