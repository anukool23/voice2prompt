//go:build darwin

#import <Cocoa/Cocoa.h>

// Exported Go callbacks.
extern void goTrayOpen(void);
extern void goTrayToggle(void);
extern void goTrayQuit(void);

@interface V2PTrayHandler : NSObject
@end

@implementation V2PTrayHandler
- (void)open:(id)sender { goTrayOpen(); }
- (void)toggle:(id)sender { goTrayToggle(); }
- (void)quit:(id)sender { goTrayQuit(); }
@end

static NSStatusItem *gStatusItem;
static V2PTrayHandler *gHandler;

void trayStart(const char *iconPath) {
    // Create UI on the main queue so it works alongside the Wails run loop.
    dispatch_async(dispatch_get_main_queue(), ^{
        gHandler = [[V2PTrayHandler alloc] init];
        gStatusItem = [[NSStatusBar systemStatusBar]
            statusItemWithLength:NSVariableStatusItemLength];

        NSImage *img = [[NSImage alloc]
            initWithContentsOfFile:[NSString stringWithUTF8String:iconPath]];
        if (img != nil) {
            [img setSize:NSMakeSize(18, 18)];
            [img setTemplate:YES]; // adapt to light/dark menu bar
            gStatusItem.button.image = img;
        } else {
            gStatusItem.button.title = @"V2P";
        }

        NSMenu *menu = [[NSMenu alloc] init];

        NSMenuItem *openItem = [[NSMenuItem alloc]
            initWithTitle:@"Open Voice2Prompt" action:@selector(open:) keyEquivalent:@""];
        openItem.target = gHandler;
        [menu addItem:openItem];

        NSMenuItem *toggleItem = [[NSMenuItem alloc]
            initWithTitle:@"Start / Stop Dictation" action:@selector(toggle:) keyEquivalent:@""];
        toggleItem.target = gHandler;
        [menu addItem:toggleItem];

        [menu addItem:[NSMenuItem separatorItem]];

        NSMenuItem *quitItem = [[NSMenuItem alloc]
            initWithTitle:@"Quit Voice2Prompt" action:@selector(quit:) keyEquivalent:@"q"];
        quitItem.target = gHandler;
        [menu addItem:quitItem];

        gStatusItem.menu = menu;
    });
}
