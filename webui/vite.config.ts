import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            '@': path.resolve(__dirname, 'src'),
        },
    },
    server: {
        host: true,
        allowedHosts: ['dev-docker.voxon.lan'],
        proxy: {
            '/pricer': {
                target: 'http://pricer:8080',
                changeOrigin: true,
                rewrite: path => path.replace(/^\/pricer/, ''),
            },
            '/researcher': {
                target: 'http://researcher:8080',
                changeOrigin: true,
                rewrite: path => path.replace(/^\/researcher/, ''),
            },
        },
    },
})
