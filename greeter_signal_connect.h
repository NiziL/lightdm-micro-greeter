extern void authentication_complete_cb(LightDMGreeter *greeter);
extern void show_prompt_cb(LightDMGreeter *greeter, char *text, LightDMPromptType type);

void greeter_signal_connect(LightDMGreeter* greeter);
