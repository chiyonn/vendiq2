import { Routes, Route, Navigate } from 'react-router-dom';
import PricingPage from "@/pages/PricingPage";
import PricingQueuesPage from "@/pages/PricingQueuesPage";
import MainLayout from "@/layouts/MainLayout";
import "./App.css";

function App() {
    return (
        <Routes>
            <Route element={<MainLayout />}>
                <Route path="/" element={<Navigate to="/pricings" replace />} />
                <Route path="/pricings" element={<PricingPage />} />
                <Route path="/pricings/queues" element={<PricingQueuesPage />} />
            </Route>
        </Routes>
    );
}

export default App;
