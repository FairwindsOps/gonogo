// To see all options:
// https://vuepress.vuejs.org/config/
// https://vuepress.vuejs.org/theme/default-theme-config.html
module.exports = {
    title: "PROJECT-NAME Documentation",
    description: "Documentation for Fairwinds' GoNoGo",
    themeConfig: {
        docsRepo: "FairwindsOps/GoNoGo",
        sidebar: [
            {
                title: "GoNoGo",
                path: "/",
                sidebarDepth: 0,
            },
            {
                title: "Installation",
                path: "/installation",
            },
            {
                title: "Quickstart",
                path: "/quickstart",
            },
            {
                title: "FAQ",
                path: "/faq",
            }
        ]

    },
}