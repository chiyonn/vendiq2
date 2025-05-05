import { Routes, Route } from 'react-router-dom';
import PricingPage from "@/pages/PricingPage";
import MainLayout from "@/layouts/MainLayout";
import "./App.css";

function App() {
    return (
        <Routes>
            <Route element={<MainLayout />}>
                <Route path="/" element={<PricingPage />} />
            </Route>
        </Routes>
    );
}

export default App;
