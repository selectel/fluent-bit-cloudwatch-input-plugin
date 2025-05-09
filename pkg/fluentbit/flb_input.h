/* -*- Mode: C; tab-width: 4; indent-tabs-mode: nil; c-basic-offset: 4 -*- */

/*  Fluent Bit Go!
 *  ==============
 *  Copyright (C) 2022 The Fluent Bit Go Authors
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

#ifndef FLBGO_INPUT_H
#define FLBGO_INPUT_H

struct flb_api {
    char *_;
    char *(*input_get_property) (char *, void *);
};

struct flb_plugin_input_proxy_context {
    int coll_fd;
   /* This context is set by the remote init and is passed to remote collect */
   void *remote_context;
   /* A proxy ptr is needed to store the proxy type/lang (OUTPUT/GOLANG) */
   struct flb_plugin_proxy *proxy;
};

/* This structure is used for initialization.
 * It matches the one in proxy/go/go.c in fluent-bit source code.
 */
struct flbgo_input_plugin {
    void *_;
    struct flb_api *api;
    struct flb_input_instance *i_ins;
    struct flb_plugin_input_proxy_context *context;
};

char *input_get_property(char *key, void *plugin)
{
    struct flbgo_input_plugin *p = plugin;
    return p->api->input_get_property(key, p->i_ins);
}

#endif
