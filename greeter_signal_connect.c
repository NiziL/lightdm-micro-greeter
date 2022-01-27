#include "_cgo_export.h"

void greeter_signal_connect(LightDMGreeter* greeter) {
    g_signal_connect(greeter, 
                     "authentication-complete", 
                     G_CALLBACK(authentication_complete_cb), 
                     NULL);

    g_signal_connect(greeter, 
                     "show-prompt", 
                     G_CALLBACK(show_prompt_cb), 
                     NULL);
}
