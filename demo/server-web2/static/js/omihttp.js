class OmiHttp {
    constructor(baseUrl) {
        this.baseUrl = baseUrl || ''; // 基础 URL，默认为空字符串
    }

    // 超时封装
    async fetchWithTimeout(resource, options = {}, timeout = 5000) {
        const controller = new AbortController();
        const id = setTimeout(() => controller.abort(), timeout);
        const response = await fetch(resource, { ...options, signal: controller.signal });
        clearTimeout(id);
        return response;
    }

    // GET 请求
    async get(path) {
        try {
            const response = await this.fetchWithTimeout(`${this.baseUrl}${path}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                throw new Error(`GET request failed with status ${response.status}: ${await response.text()}`);
            }

            return await response.json();
        } catch (error) {
            console.error('Error in GET request:', error);
            return { error: error.message };
        }
    }

    // POST 请求
    async post(path, value) {
        try {
            const response = await this.fetchWithTimeout(`${this.baseUrl}${path}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(value)
            });

            if (!response.ok) {
                throw new Error(`POST request failed with status ${response.status}: ${await response.text()}`);
            }

            return await response.json();
        } catch (error) {
            console.error('Error in POST request:', error);
            return { error: error.message };
        }
    }
}

// 实例化并导出 omihttp
const omihttp = new OmiHttp(); // 替换为实际的基础 URL
export default omihttp;