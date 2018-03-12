// Objective-C API for talking to github.com/textileio/textile-go/mobile Go package.
//   gobind -lang=objc github.com/textileio/textile-go/mobile
//
// File is generated by gobind. Do not edit.

#ifndef __Mobile_H__
#define __Mobile_H__

@import Foundation;
#include "Universe.objc.h"


@class MobileMobile;
@class MobileMobileConfig;
@class MobileNode;

@interface MobileMobile : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) id _ref;

- (instancetype)initWithRef:(id)ref;
- (instancetype)init;
// skipped method Mobile.NewNode with unsupported parameter or return types

@end

@interface MobileMobileConfig : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) id _ref;

- (instancetype)initWithRef:(id)ref;
- (instancetype)init;
/**
 * Path for the node's data directory
 */
- (NSString*)repoPath;
- (void)setRepoPath:(NSString*)v;
@end

@interface MobileNode : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) id _ref;

- (instancetype)initWithRef:(id)ref;
- (instancetype)init;
- (BOOL)start:(NSError**)error;
- (BOOL)stop:(NSError**)error;
@end

FOUNDATION_EXPORT MobileNode* MobileNewTextile(NSString* repoPath);

#endif
