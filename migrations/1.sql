/*
 * Copyright 2020 The Ledger Authors
 *
 * Licensed under the AGPL, Version 3.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.gnu.org/licenses/agpl-3.0.en.html
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
create table if not exists secret_tokens
(
    data text not null
        constraint secret_tokens_data_key
            unique
        constraint no_dash_character
            check (strpos(data, '-'::text) = 0)
);

-- alter table secret_tokens owner to interviewed;
