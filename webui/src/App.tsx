import { Routes, Route, Navigate } from 'react-router-dom';
import PricingPage from "@/pages/PricingPage";
import PricingQueuesPage from "@/pages/PricingQueuesPage";
import MainLayout from "@/layouts/MainLayout";
import "./App.css";

function App() {
    return (
        <Routes>
            <Route element={<MainLayout />}>
                <Route path="/" element={<Navigate to="/pricing" replace />} />
                <Route path="/pricing" element={<PricingPage />} />
                <Route path="/pricing/queues" element={<PricingQueuesPage />} />
            </Route>
        </Routes>
    );
}

export default App;
