# SPDX-License-Identifier: MIT
# FileENIAC Third-Party Licenses

This document lists all third-party open-source components used by FileENIAC
and their respective licenses.

## Go Dependencies (Backend)

| Package | License | Version | Copyright |
|---------|---------|---------|-----------|
| github.com/BurntSushi/toml | MIT | v1.3.2 | Copyright (c) The BurntSushi contributors |
| github.com/google/uuid | BSD-3-Clause | v1.6.0 | Copyright (c) The Go Authors |
| github.com/jlaffaye/ftp | BSD-3-Clause | v0.2.1 | Copyright (c) Julien Laffaye |
| github.com/mattn/go-sqlite3 | BSD-3-Clause | v1.14.22 | Copyright (c) The Authors |
| github.com/spf13/cobra | Apache-2.0 | v1.8.1 | Copyright 2013-2024 The Cobra Authors |
| go.uber.org/zap | MIT | v1.27.0 | Copyright (c) Uber Technologies, Inc. |
| go.uber.org/multierr | MIT | v1.10.0 | Copyright (c) Uber Technologies, Inc. |
| github.com/inconshreveable/mousetrap | Apache-2.0 | v1.1.0 | Copyright (c) 2016 Inconshreveable |
| github.com/spf13/pflag | BSD-3-Clause | v1.0.5 | Copyright (c) 2012 The Go Authors |

---

## Rust Dependencies (Desktop — Tauri)

| Crate | License | Version |
|-------|---------|---------|
| tauri | MIT/Apache-2.0 | 2.x |
| tauri-plugin-opener | MIT | 2.x |
| tauri-plugin-dialog | MIT | 2.x |
| serde | MIT/Apache-2.0 | 1.x |
| serde_json | MIT/Apache-2.0 | 1.x |
| windows-sys | MIT/Apache-2.0 | (Windows bindings) |
| tokio | MIT | async runtime |
| tracing | MIT | logging |
| uuid | Apache-2.0/MIT | v1.x |

Note: This is a non-exhaustive list of notable crates. The full dependency
tree contains many additional crates, all of which are open-source under
MIT, Apache-2.0, BSD-2, BSD-3, or similar permissive licenses. See
`apps/desktop/src-tauri/Cargo.lock` for the complete dependency list with
license information.

---

## Node.js Dependencies (Frontend)

| Package | License | Version |
|---------|---------|---------|
| react | MIT | 18.3.1 |
| react-dom | MIT | 18.3.1 |
| react-router-dom | MIT | 6.26.0 |
| @tauri-apps/api | MIT | 2.x |
| @tauri-apps/plugin-dialog | MIT | 2.7.1 |
| vite | MIT | 5.3.4 |
| tailwindcss | MIT | 3.4.4 |
| typescript | Apache-2.0/BSD-2 | 5.5.3 |
| vitest | MIT | 2.x |
| @testing-library/react | MIT | 16.x |
| autoprefixer | MIT | 10.4.19 |
| postcss | MIT | 8.4.39 |
| @vitejs/plugin-react | MIT | 4.3.1 |
| sharp | Apache-2.0 | 0.35.1 |
| jsdom | MIT | 24.x |

Note: This is a non-exhaustive list of direct dependencies. The full
dependency tree contains many additional packages. All are open-source
under permissive licenses. See `pnpm-lock.yaml` for the complete list.

---

## License Texts

### MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

### BSD-3-Clause

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software without
   specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

### Apache License 2.0

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.

---

## Compliance

FileENIAC is committed to open-source compliance. All third-party components
used in this project are open-source software with permissive licenses that
allow redistribution under the MIT License terms.

If you believe any third-party attribution is missing or incorrect, please
open an issue on the repository.