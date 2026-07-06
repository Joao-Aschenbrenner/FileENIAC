// SPDX-License-Identifier: MIT
import { Routes, Route } from "react-router-dom";
import ErrorBoundary from "./components/ErrorBoundary";
import Layout from "./components/Layout";
import Onboarding from "./pages/Onboarding";
import Dashboard from "./pages/Dashboard";
import Projects from "./pages/Projects";
import Servers from "./pages/Servers";
import History from "./pages/History";
import ProjectDetails from "./pages/ProjectDetails";
import DeployCenter from "./pages/DeployCenter";
import RollbackCenter from "./pages/RollbackCenter";
import DiffViewer from "./pages/DiffViewer";
import SyncCenter from "./pages/SyncCenter";
import HealthCenter from "./pages/HealthCenter";
import GitHubLogin from "./pages/GitHubLogin";
import GitHubOrgs from "./pages/GitHubOrgs";
import GitHubRepos from "./pages/GitHubRepos";
import WorkspaceBootstrap from "./pages/WorkspaceBootstrap";

export default function App() {
  return (
    <ErrorBoundary>
      <Routes>
        <Route path="/" element={<Onboarding />} />
        <Route element={<Layout />}>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/configurar" element={<WorkspaceBootstrap />} />
          <Route path="/bootstrap" element={<WorkspaceBootstrap />} />
          <Route path="/projects" element={<Projects />} />
          <Route path="/projects/:name" element={<ProjectDetails />} />
          <Route path="/servers" element={<Servers />} />
          <Route path="/history" element={<History />} />
          <Route path="/deploy" element={<DeployCenter />} />
          <Route path="/rollback" element={<RollbackCenter />} />
          <Route path="/diff" element={<DiffViewer />} />
          <Route path="/sync" element={<SyncCenter />} />
          <Route path="/health" element={<HealthCenter />} />
          <Route path="/github/login" element={<GitHubLogin />} />
          <Route path="/github/orgs" element={<GitHubOrgs />} />
          <Route path="/github/repos" element={<GitHubRepos />} />
        </Route>
      </Routes>
    </ErrorBoundary>
  );
}