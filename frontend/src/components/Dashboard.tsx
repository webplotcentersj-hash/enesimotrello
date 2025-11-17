import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { boardAPI, taskAPI } from '../services/api.ts';
import { useAuth } from '../services/AuthContext.tsx';

interface Board {
  id: number;
  title: string;
  description: string;
  created_at: string;
}

const Dashboard: React.FC = () => {
  const [boards, setBoards] = useState<Board[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [newBoard, setNewBoard] = useState({ title: '', description: '' });
  const [totalTasks, setTotalTasks] = useState(0);
  const [completedTasks, setCompletedTasks] = useState(0);
  const [activeTasks, setActiveTasks] = useState(0);
  const { displayName, anonymousUserId, loading: authLoading } = useAuth();

  useEffect(() => {
    if (!authLoading && anonymousUserId) {
      fetchBoards();
      fetchTaskStats();
    }
  }, [authLoading, anonymousUserId]);

  const fetchBoards = async () => {
    try {
      if (authLoading || !anonymousUserId) return;
      const response = await boardAPI.getBoards();
      setBoards(response.data.boards || []);
    } catch (error) {
      console.error('Failed to fetch boards:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchTaskStats = async () => {
    try {
      if (authLoading || !anonymousUserId) return;
      const response = await boardAPI.getBoards();
      const allBoards = response.data.boards || [];
      
      let total = 0;
      let completed = 0;
      let active = 0;

      // Fetch tasks for each board
      for (const board of allBoards) {
        try {
          const tasksResponse = await taskAPI.getTasks(board.id);
          const tasks = tasksResponse.data.tasks || [];
          total += tasks.length;
          completed += tasks.filter((t: any) => t.status === 'done').length;
          active += tasks.filter((t: any) => t.status === 'in_progress').length;
        } catch (err) {
          console.error(`Failed to fetch tasks for board ${board.id}:`, err);
        }
      }

      setTotalTasks(total);
      setCompletedTasks(completed);
      setActiveTasks(active);
    } catch (error) {
      console.error('Failed to fetch task stats:', error);
    }
  };

  const handleCreateBoard = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (authLoading || !anonymousUserId) return;
      await boardAPI.createBoard(newBoard);
      setNewBoard({ title: '', description: '' });
      setShowModal(false);
      fetchBoards();
      fetchTaskStats();
    } catch (error) {
      console.error('Failed to create board:', error);
    }
  };

  const handleDeleteBoard = async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this board?')) return;
    try {
      await boardAPI.deleteBoard(id);
      fetchBoards();
      fetchTaskStats();
    } catch (error) {
      console.error('Failed to delete board:', error);
    }
  };

  const getProductivity = () => {
    if (totalTasks === 0) return 0;
    return Math.round((completedTasks / totalTasks) * 100);
  };

  const getColorClass = (index: number) => {
    const colors = [
      'from-orange-200 to-orange-300',
      'from-orange-300 to-orange-400',
      'from-orange-200 to-orange-400',
      'from-orange-100 to-orange-300',
      'from-orange-300 to-orange-500',
      'from-orange-200 to-orange-500',
    ];
    return colors[index % colors.length];
  };

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg text-gray-600">Loading your workspace...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      {/* Header Section */}
      <div className="mb-8">
        <h1 className="text-4xl font-bold text-gray-800 mb-2">
          Welcome back, {displayName}! 👋
        </h1>
        <p className="text-gray-600">Here are your boards and projects</p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="bg-gradient-to-br from-orange-200 to-orange-300 rounded-2xl p-6 text-white shadow-lg transform hover:scale-105 transition-all duration-200">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-blue-100 text-sm font-medium">Total Boards</p>
              <p className="text-3xl font-bold mt-2">{boards.length}</p>
            </div>
            <div className="text-5xl opacity-50">📊</div>
          </div>
        </div>

        <div className="bg-gradient-to-br from-orange-300 to-orange-400 rounded-2xl p-6 text-white shadow-lg transform hover:scale-105 transition-all duration-200">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-purple-100 text-sm font-medium">Active Tasks</p>
              <p className="text-3xl font-bold mt-2">{activeTasks}</p>
            </div>
            <div className="text-5xl opacity-50">⚡</div>
          </div>
        </div>

        <div className="bg-gradient-to-br from-orange-400 to-orange-500 rounded-2xl p-6 text-white shadow-lg transform hover:scale-105 transition-all duration-200">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-pink-100 text-sm font-medium">Completed</p>
              <p className="text-3xl font-bold mt-2">{completedTasks}</p>
            </div>
            <div className="text-5xl opacity-50">🎉</div>
          </div>
        </div>

        <div className="bg-gradient-to-br from-orange-200 to-orange-400 rounded-2xl p-6 text-white shadow-lg transform hover:scale-105 transition-all duration-200">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-green-100 text-sm font-medium">Productivity</p>
              <p className="text-3xl font-bold mt-2">{getProductivity()}%</p>
              <p className="text-xs text-green-100 mt-1">{completedTasks}/{totalTasks} tasks</p>
            </div>
            <div className="text-5xl opacity-50">📈</div>
          </div>
        </div>
      </div>

      {/* Boards Section */}
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Your Boards</h2>
        <button
          onClick={() => setShowModal(true)}
          className="flex items-center space-x-2 bg-gradient-to-r from-orange-300 to-orange-400 hover:from-orange-400 hover:to-orange-500 text-white px-6 py-3 rounded-xl shadow-lg transform hover:scale-105 transition-all duration-200 font-semibold"
        >
          <span className="text-xl">+</span>
          <span>New Board</span>
        </button>
      </div>

      {loading ? (
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-indigo-600"></div>
        </div>
      ) : boards.length === 0 ? (
        <div className="text-center py-20">
          <div className="text-8xl mb-4">📋</div>
          <h3 className="text-2xl font-bold text-gray-700 mb-2">No boards yet</h3>
          <p className="text-gray-500 mb-6">Create your first board to get started!</p>
          <button
            onClick={() => setShowModal(true)}
            className="bg-gradient-to-r from-orange-300 to-orange-400 hover:from-orange-400 hover:to-orange-500 text-white px-8 py-3 rounded-xl shadow-lg transform hover:scale-105 transition-all duration-200 font-semibold"
          >
            Create Your First Board
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {boards.map((board, index) => (
            <div
              key={board.id}
              className="bg-white rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-2 overflow-hidden"
            >
              <div className={`h-32 bg-gradient-to-br ${getColorClass(index)} p-6 relative overflow-hidden`}>
                <div className="absolute top-0 right-0 w-32 h-32 bg-white opacity-10 rounded-full -mr-16 -mt-16"></div>
                <div className="absolute bottom-0 left-0 w-24 h-24 bg-white opacity-10 rounded-full -ml-12 -mb-12"></div>
                <h3 className="text-2xl font-bold text-white relative z-10">{board.title}</h3>
              </div>
              <div className="p-6">
                <p className="text-gray-600 mb-4 h-12 overflow-hidden">
                  {board.description || 'No description'}
                </p>
                <div className="flex items-center justify-between text-sm text-gray-500 mb-4">
                  <span>📅 {new Date(board.created_at).toLocaleDateString()}</span>
                  <span>0 tasks</span>
                </div>
                <div className="flex space-x-2">
                  <Link
                    to={`/board/${board.id}`}
                    className="flex-1 bg-indigo-50 hover:bg-indigo-100 text-indigo-600 font-semibold py-2 px-4 rounded-lg transition-colors text-center"
                  >
                    Open Board
                  </Link>
                  <button
                    onClick={() => handleDeleteBoard(board.id)}
                    className="bg-red-50 hover:bg-red-100 text-red-600 font-semibold py-2 px-4 rounded-lg transition-colors"
                  >
                    🗑️
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create Board Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full p-8 transform transition-all">
            <div className="flex justify-between items-center mb-6">
              <h3 className="text-2xl font-bold text-gray-800">Create New Board</h3>
              <button
                onClick={() => setShowModal(false)}
                className="text-gray-400 hover:text-gray-600 text-2xl"
              >
                ×
              </button>
            </div>
            <form onSubmit={handleCreateBoard} className="space-y-4">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Board Title
                </label>
                <input
                  type="text"
                  required
                  value={newBoard.title}
                  onChange={(e) => setNewBoard({ ...newBoard, title: e.target.value })}
                  className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all outline-none"
                  placeholder="My Awesome Project"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Description
                </label>
                <textarea
                  value={newBoard.description}
                  onChange={(e) => setNewBoard({ ...newBoard, description: e.target.value })}
                  className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all outline-none resize-none"
                  rows={3}
                  placeholder="What's this board about?"
                />
              </div>
              <div className="flex space-x-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-700 font-semibold py-3 px-4 rounded-xl transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="flex-1 bg-gradient-to-r from-orange-300 to-orange-400 hover:from-orange-400 hover:to-orange-500 text-white font-semibold py-3 px-4 rounded-xl shadow-lg transform hover:scale-105 transition-all duration-200"
                >
                  Create Board
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;
