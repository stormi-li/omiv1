const { ref, reactive, computed, onMounted, createApp } = Vue;

const App = {
    setup() {
        // 状态变量
        const servers = ref([]); // 全部服务器数据
        const searchQuery = ref(""); // 搜索关键词
        const view = ref("list"); // 当前视图 ("list" or "details")
        const selectedServer = reactive({
            name: "",
            address: "",
            weight: 0,
            details: {},
            commands: []
        });

        // 过滤后的服务器
        const filteredServers = computed(() =>
            servers.value.filter(server =>
                server.name.toLowerCase().includes(searchQuery.value.toLowerCase())
            )
        );

        // 加载服务器数据 (模拟 API 请求)
        const loadServers = async () => {
            servers.value = await new Promise(resolve => {
                setTimeout(() => {
                    resolve([
                        {
                            id: 1,
                            name: "Login-Server",
                            address: "localhost:8080",
                            weight: 1
                        },
                        {
                            id: 2,
                            name: "Data-Server",
                            address: "localhost:9090",
                            weight: 2
                        }
                    ]);
                }, 500);
            });
        };

        // 切换到详情视图
        const showDetails = server => {
            Object.assign(selectedServer, {
                ...server,
                details: {
                    weight: server.weight,
                    start_time: "2024-11-27 02:33:26",
                    run_time: "9m33.113991788s"
                },
                commands: [
                    { name: "update_weight" },
                    { name: "open_cache" },
                    { name: "update_cache_size" }
                ]
            });
            view.value = "details";
        };

        // 返回列表视图
        const goBack = () => {
            view.value = "list";
        };

        // 刷新详情数据
        const refreshDetails = () => {
            selectedServer.details.run_time = new Date().toLocaleTimeString();
        };

        // 生命周期钩子
        onMounted(() => {
            loadServers();
        });

        // 返回所有绑定数据和方法
        return {
            servers,
            searchQuery,
            view,
            selectedServer,
            filteredServers,
            showDetails,
            goBack,
            refreshDetails
        };
    }
};

createApp(App).mount("#app");