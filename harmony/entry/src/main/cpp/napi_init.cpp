#include "libfull.h"
#include "napi/native_api.h"
#include "string.h"
#include "thread"
#include "future"
#include "memory"

// 工具函数:从 napi_value 转换为 C 字符串
// 注意:调用者负责释放返回的内存
static char *value2String(napi_env env, napi_value value) {
    size_t len = 0;
    napi_get_value_string_utf8(env, value, nullptr, 0, &len);
    char *buf = new char[len + 1];
    napi_get_value_string_utf8(env, value, buf, len + 1, &len);
    return buf;
}

// // Add 函数的 NAPI 包装 - 同步返回结果
// static napi_value Add0(napi_env env, napi_callback_info info) {
//     napi_value result;
    
//     // 获取参数
//     size_t argc = 2;
//     napi_value args[2] = {nullptr, nullptr};
//     napi_get_cb_info(env, info, &argc, args, nullptr, nullptr);
    
//     // 转换参数类型
//     double x, y;
//     napi_get_value_double(env, args[0], &x);
//     napi_get_value_double(env, args[1], &y);
    
//     // 使用 shared_ptr 管理 promise,避免栈引用问题
//     auto promise_ptr = std::make_shared<std::promise<double>>();
//     std::future<double> future = promise_ptr->get_future();
    
//     std::thread t([promise_ptr, x, y]() {
//         try {
//             double sum = Add(x, y);
//             promise_ptr->set_value(sum);
//         } catch (...) {
//             promise_ptr->set_exception(std::current_exception());
//         }
//     });
//     t.join();
    
//     // 获取结果并返回
//     try {
//         double ret = future.get();
//         napi_create_double(env, ret, &result);
//     } catch (...) {
//         // 发生异常时返回 NaN
//         napi_create_double(env, 0.0, &result);
//     }
    
//     return result;
// }

// // Hello 函数的 NAPI 包装 - 返回字符串
// static napi_value Hello0(napi_env env, napi_callback_info info) {
//     napi_value result;
    
//     // 使用 shared_ptr 管理 promise,避免栈引用问题
//     auto promise_ptr = std::make_shared<std::promise<char*>>();
//     std::future<char*> future = promise_ptr->get_future();
    
//     std::thread t([promise_ptr]() {
//         try {
//             char* msg = Hello();
//             promise_ptr->set_value(msg);
//         } catch (...) {
//             promise_ptr->set_exception(std::current_exception());
//         }
//     });
//     t.join();
    
//     // 获取结果并转换为 NAPI 字符串
//     char* msg = nullptr;
//     try {
//         msg = future.get();
//         if (msg != nullptr) {
//             napi_create_string_utf8(env, msg, strlen(msg), &result);
//             // 释放 Go 分配的内存
//             free(msg);
//         } else {
//             napi_create_string_utf8(env, "", 0, &result);
//         }
//     } catch (...) {
//         // 发生异常时返回空字符串
//         napi_create_string_utf8(env, "", 0, &result);
//         if (msg != nullptr) {
//             free(msg);
//         }
//     }
    
//     return result;
// }

// FullSdkRun 函数的 NAPI 包装 - 接收两个字符串参数并返回整数
static napi_value FullSdkRun0(napi_env env, napi_callback_info info) {
    napi_value result;
    
    // 获取参数
    size_t argc = 2;
    napi_value args[2] = {nullptr, nullptr};
    napi_get_cb_info(env, info, &argc, args, nullptr, nullptr);
    
    // 检查参数数量
    if (argc < 2) {
        napi_create_int32(env, -1, &result);
        return result;
    }
    
    // 转换参数为 C 字符串
    char* logDirPath = value2String(env, args[0]);
    char* dataBaseDirPath = value2String(env, args[1]);
    
    // 使用 shared_ptr 管理 promise,避免栈引用问题
    auto promise_ptr = std::make_shared<std::promise<int>>();
    std::future<int> future = promise_ptr->get_future();
    
    std::thread t([promise_ptr, logDirPath, dataBaseDirPath]() {
        try {
            int ret = FullSdkRun(logDirPath, dataBaseDirPath);
            promise_ptr->set_value(ret);
        } catch (...) {
            promise_ptr->set_exception(std::current_exception());
        }
        // 在线程内释放字符串内存
        delete[] logDirPath;
        delete[] dataBaseDirPath;
    });
    t.join();
    
    // 获取结果并返回
    try {
        int ret = future.get();
        napi_create_int32(env, ret, &result);
    } catch (...) {
        // 发生异常时返回 -1
        napi_create_int32(env, -1, &result);
    }
    
    return result;
}

// 模块初始化
EXTERN_C_START
static napi_value Init(napi_env env, napi_value exports) {
    napi_property_descriptor desc[] = {
        // {"add", nullptr, Add0, nullptr, nullptr, nullptr, napi_default, nullptr},
        // {"hello", nullptr, Hello0, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"fullSdkRun", nullptr, FullSdkRun0, nullptr, nullptr, nullptr, napi_default, nullptr},
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