import { hilog } from '@kit.PerformanceAnalysisKit';
import napi from 'libentry.so';

const DOMAIN = 0x0000;

@Entry
@Component
struct Index {
  @State message: string = 'Hello World';

  build() {
    Row() {
      Column() {
        Text(this.message)
          .fontSize($r('app.float.page_text_font_size'))
          .fontWeight(FontWeight.Bold)
          .onClick(() => {
            this.message = 'Welcome';
            hilog.info(DOMAIN, 'testTag', 'Test NAPI napi.add 2 + 3 = %{public}d', napi.add(2.0, 3.0));
            hilog.info(DOMAIN, 'testTag', 'Test NAPI napi.hello %{public}s', napi.hello());
          })
      }
      .width('100%')
    }
    .height('100%')
  }
}
