import { Outlet } from 'react-router-dom';
import PageNav from "@/components/global/PageNav";
import styles from './MainLayout.module.css';

const MainLayout = () => {
    return (
        <div className={styles.container}>
            <header className={styles.header}>ヘッダー</header>
            <div className={styles.contentWrapper}>
                <aside className={styles.sidebarLeft}><PageNav /></aside>
                <main className={styles.main}>
                    <Outlet />
                </main>
                <aside className={styles.sidebarRight}>右サイドバー</aside>
            </div>
            <footer className={styles.footer}>フッター</footer>
        </div>
    );
};

export default MainLayout;
