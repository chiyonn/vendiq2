import { Routes, Route, Navigate } from 'react-router-dom';
import PricingPage from "@/pages/PricingPage";
import MainLayout from "@/layouts/MainLayout";
import "./App.css";

function App() {
    return (
        <Routes>
            <Route element={<MainLayout />}>
                <Route path="/" element={<Navigate to="/pricing" replace />} />
                <Route path="/pricing" element={<PricingPage />} />
            </Route>
        </Routes>
    );
}

export default App;
