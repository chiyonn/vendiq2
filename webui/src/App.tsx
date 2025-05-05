import { Routes, Route } from 'react-router-dom';
import PricingPage from "@/pages/PricingPage";

function App() {
    return (
        <Routes>
            <Route path="/pricing" element={<PricingPage />} />
        </Routes>
    );
}

export default App;
