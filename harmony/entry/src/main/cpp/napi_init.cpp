#include "libfull.h"
#include "napi/native_api.h"
#include "string.h"
#include "thread"
#include "future"

// 工具函数:从 napi_value 转换为 C 字符串
static char *value2String(napi_env env, napi_value value) {
    size_t len = 0;
    napi_get_value_string_utf8(env, value, nullptr, 0, &len);
    char *buf = new char[len + 1];
    napi_get_value_string_utf8(env, value, buf, len + 1, &len);
    return buf;
}

// Add 函数的 NAPI 包装 - 同步返回结果
static napi_value Add0(napi_env env, napi_callback_info info) {
    napi_value result;
    
    // 获取参数
    size_t argc = 2;
    napi_value args[2] = {nullptr, nullptr};
    napi_get_cb_info(env, info, &argc, args, nullptr, nullptr);
    
    // 转换参数类型
    double x, y;
    napi_get_value_double(env, args[0], &x);
    napi_get_value_double(env, args[1], &y);
    
    // 使用 promise/future 在子线程中执行 Go 函数
    std::promise<double> promise;
    std::future<double> future = promise.get_future();
    
    std::thread t([&promise, x, y]() {
        double sum = Add(x, y);
        promise.set_value(sum);
    });
    t.join();
    
    // 获取结果并返回
    double ret = future.get();
    napi_create_double(env, ret, &result);
    return result;
}

// Hello 函数的 NAPI 包装 - 返回字符串
static napi_value Hello0(napi_env env, napi_callback_info info) {
    napi_value result;
    
    // 使用 promise/future 在子线程中执行 Go 函数
    std::promise<char*> promise;
    std::future<char*> future = promise.get_future();
    
    std::thread t([&promise]() {
        char* msg = Hello();
        promise.set_value(msg);
    });
    t.join();
    
    // 获取结果并转换为 NAPI 字符串
    char* msg = future.get();
    napi_create_string_utf8(env, msg, strlen(msg), &result);
    return result;
}

// 模块初始化
EXTERN_C_START
static napi_value Init(napi_env env, napi_value exports) {
    napi_property_descriptor desc[] = {
        {"add", nullptr, Add0, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"hello", nullptr, Hello0, nullptr, nullptr, nullptr, napi_default, nullptr},
    };
    napi_define_properties(env, exports, sizeof(desc) / sizeof(desc[0]), desc);
    return exports;
}
EXTERN_C_END

// 模块注册
static napi_module demoModule = {
    .nm_version = 1,
    .nm_flags = 0,
    .nm_filename = nullptr,
    .nm_register_func = Init,
    .nm_modname = "entry",  // 模块名称
    .nm_priv = ((void *)0),
    .reserved = {0},
};

extern "C" __attribute__((constructor)) void RegisterHarmonyModule(void) {
    napi_module_register(&demoModule);
}