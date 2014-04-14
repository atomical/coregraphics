// +build darwin

package coregraphics

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa
// #include <QuartzCore/QuartzCore.h>
//
// void * DictionaryObjectForKey(NSDictionary *dict, id key ){
//   return [dict objectForKey: key];
// }
//
// static NSString* GetDictValue( NSDictionary* dict, CFStringRef key )
// {
//     if ([dict objectForKey:(id)kCGWindowOwnerPID])
//     {
//        NSString *string = 
//          [NSString stringWithFormat:@"(%@)", [dict objectForKey:(id)key]];
//        return string;
//      }
//     else 
//     {
//         return nil;
//     }
// }
//
// const char * NSStringToCString( NSString * str ){
//   const char * c_str = [str UTF8String];
//   return c_str;
// }
//
import "C"
import "unsafe"

// Quartz Window Services Reference
// https://developer.apple.com/library/mac/documentation/Carbon/reference/CGWindow_Reference/Reference/Functions.html

// CGWindowID

const (
  KCGNullWindowID           = C.kCGNullWindowID
  KCGWindowSharingNone      = C.kCGWindowSharingNone
  KCGWindowSharingReadOnly  = C.kCGWindowSharingReadOnly
  KCGWindowSharingReadWrite = C.kCGWindowSharingReadWrite
)

// Window List Option Constants
// https://developer.apple.com/library/mac/documentation/Carbon/reference/CGWindow_Reference/Constants/Constants.html#//apple_ref/doc/constant_group/Window_List_Option_Constants

const (
  KCGWindowListOptionAll                 = C.kCGWindowListOptionAll
  KCGWindowListOptionOnScreenOnly        = C.kCGWindowListOptionOnScreenOnly
  KCGWindowListOptionOnScreenAboveWindow = C.kCGWindowListOptionOnScreenAboveWindow
  KCGWindowListOptionOnScreenBelowWindow = C.kCGWindowListOptionOnScreenBelowWindow
  KCGWindowListOptionIncludingWindow     = C.kCGWindowListOptionIncludingWindow
  KCGWindowListExcludeDesktopElements    = C.kCGWindowListExcludeDesktopElements
)

// const (
//   // KCGWindowOwnerName  = C.kCGWindowOwnerName
//   // KCGWindowWorkspace  = C.kCGWindowWorkspace
//   // CGWindowOwnerPID    = C.kCGWindowOwnerPID

//   // Required Window Keys
//   // const CFStringRef kCGWindowNumber;
//   // const CFStringRef kCGWindowStoreType;
//   // const CFStringRef kCGWindowLayer;
//   // const CFStringRef kCGWindowBounds;
//   // const CFStringRef kCGWindowSharingState;
//   // const CFStringRef kCGWindowAlpha;
//   // const CFStringRef kCGWindowOwnerPID;
//   // const CFStringRef kCGWindowMemoryUsage;
// )

type Rect struct {
  X float64
  Y float64
  Width float64
  Height float64
}

type Window struct {
  OwnerName string
  WindowId  int
  Rect Rect
}


func CGWindowListCopyWindowInfo( option C.CGWindowListOption, relativeToWindow C.CGWindowID) ( []Window ) {

  listArray  := C.CGWindowListCopyWindowInfo( option, relativeToWindow )
  count      := CFArrayGetCount( listArray ) 
  windows    := make([]Window, count)

  for iter := 0; iter < count; iter++ {
    entry := C.CFArrayGetValueAtIndex(listArray, C.CFIndex(iter))
    name := CFStringGet(CFDictionaryGetValue(entry, unsafe.Pointer(C.kCGWindowOwnerName)))
    windowId := CFNumberGetValue(
                  C.CFNumberRef(CFDictionaryGetValue(entry, unsafe.Pointer(C.kCGWindowNumber))),
                  C.kCGWindowIDCFNumberType,
                )
    bounds_value := CFDictionaryGetValue(entry, unsafe.Pointer(C.kCGWindowBounds))
    rect := CGRectMakeWithDictionaryRepresentation(bounds_value)

    windows[iter] = Window{
                      OwnerName: name, 
                      WindowId: windowId,
                      Rect: rect,
                    }

    // const CFStringRef kCGWindowNumber;
    // const CFStringRef kCGWindowStoreType;
    // const CFStringRef kCGWindowLayer;
    // const CFStringRef kCGWindowBounds;
    // const CFStringRef kCGWindowSharingState;
    // const CFStringRef kCGWindowAlpha;
    // const CFStringRef kCGWindowOwnerPID;
    // const CFStringRef kCGWindowMemoryUsage;

  }

  return windows
}

func CGRectMakeWithDictionaryRepresentation( dict unsafe.Pointer ) Rect {
  var rect Rect
  var CGRect C.CGRect
  
  C.CGRectMakeWithDictionaryRepresentation((*[0]byte)(dict), &CGRect)

  rect.X      = (float64)(C.CGRectGetMinX(   CGRect ))
  rect.Y      = (float64)(C.CGRectGetMinY(   CGRect ))
  rect.Width  = (float64)(C.CGRectGetWidth(  CGRect ))
  rect.Height = (float64)(C.CGRectGetHeight( CGRect ))

  return rect
}


func CFDictionaryContainsKey( entry unsafe.Pointer, key unsafe.Pointer ) int {
  return (int(C.CFDictionaryContainsKey((*[0]byte)(entry),
    unsafe.Pointer(C.kCGWindowOwnerName),
  )))
}

func CFDictionaryGetValue( entry unsafe.Pointer, key unsafe.Pointer ) unsafe.Pointer {
  return C.CFDictionaryGetValue( (*[0]byte)(entry), key )
}

func CFNumberGetValue( number C.CFNumberRef, theType C.CFNumberType ) int {
  var val int
  C.CFNumberGetValue ( number, theType, unsafe.Pointer(&val) )
  return val
}

func CFStringGet( ptr unsafe.Pointer ) string {
  length := C.CFStringGetLength((*[0]byte)(ptr))
  buffer := C.malloc(C.size_t(length + 1))

  C.CFStringGetCString( 
    (*[0]byte)(ptr), 
    (*C.char)(unsafe.Pointer(buffer)),
    length + 1,
    C.kCFStringEncodingUTF8 )

  str := C.GoStringN((*C.char)(buffer), C.int(length))
  defer C.free(buffer)

  return str
}

// CFArray Reference
// https://developer.apple.com/library/mac/documentation/corefoundation/Reference/CFArrayRef/Reference/reference.html



func CFArrayGetCount( theArray C.CFArrayRef ) ( int ) {
  return int(C.CFArrayGetCount( theArray ))
}

func CFArrayGetValueAtIndex( theArray C.CFArrayRef, idx C.CFIndex ) ( unsafe.Pointer ) {
  return C.CFArrayGetValueAtIndex(theArray, idx)
}

