//go:build darwin

#import <AVFoundation/AVFoundation.h>

// Returns the AVAuthorizationStatus for audio:
// 0 = notDetermined, 1 = restricted, 2 = denied, 3 = authorized.
int mic_status(void) {
    return (int)[AVCaptureDevice authorizationStatusForMediaType:AVMediaTypeAudio];
}

// Triggers the microphone permission prompt (no-op if already decided).
void mic_request(void) {
    [AVCaptureDevice requestAccessForMediaType:AVMediaTypeAudio
                             completionHandler:^(BOOL granted){ (void)granted; }];
}
