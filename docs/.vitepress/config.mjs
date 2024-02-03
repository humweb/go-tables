import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Go Tables",
  description: "Companion backend to inertia vue tables frontend",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Docs', link: '/introduction' }
    ],

    sidebar: [
      {
        text: 'Getting started',
        items: [
          { text: 'Introduction', link: '/introduction' },
          { text: 'Quick Start', link: '/quick-start' }
        ]
      },
      {
        text: 'Digging Deeper',
        items: [
          { text: 'Fields', link: '/fields' },
          { text: 'Filters', link: '/filters' },
          { text: 'Custom Filters', link: '/filters-custom' },
          { text: 'Relationships', link: '/relationships' },

        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/humweb/go-tables' }
    ]
  }
})
