import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { taskAPI } from '../services/api.ts';
import { resolveWebSocketUrl, useWebSocket } from '../hooks/useWebSocket.ts';

interface Task {
  id: number;
  title: string;
  description: string;
  status: 'todo' | 'in_progress' | 'done';
  priority: 'low' | 'medium' | 'high';
  assignee_id?: number;
  due_date?: string;
  created_at: string;
}

const Board: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newTask, setNewTask] = useState({
    title: '',
    description: '',
    priority: 'medium' as 'low' | 'medium' | 'high',
  });

  const websocketUrl =
    process.env.REACT_APP_WS_URL ?? resolveWebSocketUrl('');

  const { isConnected, lastMessage } = useWebSocket(websocketUrl);

  useEffect(() => {
    fetchTasks();
  }, [id]);

  useEffect(() => {
    if (lastMessage) {
      switch (lastMessage.type) {
        case 'task_created':
        case 'task_updated':
        case 'task_deleted':
          fetchTasks();
          break;
      }
    }
  }, [lastMessage]);

  const fetchTasks = async () => {
    if (!id) return;
    try {
      const response = await taskAPI.getTasks(parseInt(id));
      setTasks(response.data.tasks || []);
    } catch (error) {
      console.error('Failed to fetch tasks:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTask = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!id) return;
    try {
      await taskAPI.createTask(parseInt(id), newTask);
      setNewTask({ title: '', description: '', priority: 'medium' });
      setShowCreateModal(false);
      fetchTasks();
    } catch (error) {
      console.error('Failed to create task:', error);
    }
  };

  const handleUpdateTaskStatus = async (taskId: number, newStatus: Task['status']) => {
    if (!id) return;
    const task = tasks.find(t => t.id === taskId);
    if (!task) return;

    try {
      await taskAPI.updateTask(parseInt(id), taskId, {
        ...task,
        status: newStatus,
      });
      fetchTasks();
    } catch (error) {
      console.error('Failed to update task:', error);
    }
  };

  const getTasksByStatus = (status: Task['status']) => {
    return tasks.filter(task => task.status === status);
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high':
        return 'bg-red-100 text-red-700 border-red-200';
      case 'medium':
        return 'bg-yellow-100 text-yellow-700 border-yellow-200';
      case 'low':
        return 'bg-green-100 text-green-700 border-green-200';
      default:
        return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  };

  const getPriorityEmoji = (priority: string) => {
    switch (priority) {
      case 'high':
        return '🔴';
      case 'medium':
        return '🟡';
      case 'low':
        return '🟢';
      default:
        return '⚪';
    }
  };

  const TaskCard = ({ task }: { task: Task }) => (
    <div className="bg-white rounded-xl shadow-md hover:shadow-xl transition-all duration-300 p-4 mb-3 cursor-pointer transform hover:-translate-y-1 border-l-4 border-indigo-500">
      <div className="flex items-start justify-between mb-2">
        <h4 className="font-semibold text-gray-800 flex-1">{task.title}</h4>
        <span className={`text-xs px-2 py-1 rounded-full border ${getPriorityColor(task.priority)}`}>
          {getPriorityEmoji(task.priority)} {task.priority}
        </span>
      </div>
      {task.description && (
        <p className="text-sm text-gray-600 mb-3">{task.description}</p>
      )}
      <div className="flex items-center justify-between text-xs text-gray-500">
        <span>📅 {new Date(task.created_at).toLocaleDateString()}</span>
        <div className="flex space-x-1">
          {task.status !== 'todo' && (
            <button
              onClick={() => handleUpdateTaskStatus(task.id, 'todo')}
              className="p-1 hover:bg-gray-100 rounded"
              title="Move to To Do"
            >
              ⬅️
            </button>
          )}
          {task.status === 'todo' && (
            <button
              onClick={() => handleUpdateTaskStatus(task.id, 'in_progress')}
              className="p-1 hover:bg-blue-50 rounded text-blue-600"
              title="Start Task"
            >
              ▶️
            </button>
          )}
          {task.status === 'in_progress' && (
            <button
              onClick={() => handleUpdateTaskStatus(task.id, 'done')}
              className="p-1 hover:bg-green-50 rounded text-green-600"
              title="Complete Task"
            >
              ✅
            </button>
          )}
        </div>
      </div>
    </div>
  );

  const Column = ({ 
    title, 
    status, 
    color, 
    emoji, 
    tasks 
  }: { 
    title: string; 
    status: Task['status']; 
    color: string; 
    emoji: string; 
    tasks: Task[];
  }) => (
    <div className="bg-gray-50 rounded-2xl p-4 min-w-[320px] flex-1">
      <div className={`flex items-center justify-between mb-4 pb-3 border-b-2 ${color}`}>
        <h3 className="font-bold text-lg flex items-center space-x-2">
          <span className="text-2xl">{emoji}</span>
          <span>{title}</span>
        </h3>
        <span className="bg-white px-3 py-1 rounded-full text-sm font-semibold text-gray-700 shadow-sm">
          {tasks.length}
        </span>
      </div>
      <div className="space-y-3 max-h-[calc(100vh-300px)] overflow-y-auto pr-2 custom-scrollbar">
        {tasks.length === 0 ? (
          <div className="text-center py-8 text-gray-400">
            <p className="text-4xl mb-2">📭</p>
            <p className="text-sm">No tasks here</p>
          </div>
        ) : (
          tasks.map(task => <TaskCard key={task.id} task={task} />)
        )}
      </div>
    </div>
  );

  return (
    <div className="min-h-screen pb-8">
      {/* Header */}
      <div className="mb-8 bg-white rounded-2xl shadow-md p-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <Link 
              to="/dashboard" 
              className="text-2xl hover:bg-gray-100 p-2 rounded-lg transition-colors"
            >
              ←
            </Link>
            <div>
              <h1 className="text-3xl font-bold text-gray-800 flex items-center space-x-3">
                <span>Board #{id}</span>
                <div className="flex items-center space-x-2">
                  <div className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'}`}></div>
                  <span className="text-sm font-normal text-gray-500">
                    {isConnected ? 'Live' : 'Offline'}
                  </span>
                </div>
              </h1>
              <p className="text-gray-600">Manage your tasks efficiently</p>
            </div>
          </div>
          <button
            onClick={() => setShowCreateModal(true)}
            className="flex items-center space-x-2 bg-gradient-to-r from-orange-300 to-orange-400 hover:from-orange-400 hover:to-orange-500 text-white px-6 py-3 rounded-xl shadow-lg transform hover:scale-105 transition-all duration-200 font-semibold"
          >
            <span className="text-xl">+</span>
            <span>New Task</span>
          </button>
        </div>
      </div>

      {/* Kanban Board */}
      {loading ? (
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-indigo-600"></div>
        </div>
      ) : (
        <div className="flex space-x-6 overflow-x-auto pb-4">
          <Column
            title="To Do"
            status="todo"
            color="border-blue-500"
            emoji="📝"
            tasks={getTasksByStatus('todo')}
          />
          <Column
            title="In Progress"
            status="in_progress"
            color="border-yellow-500"
            emoji="⚡"
            tasks={getTasksByStatus('in_progress')}
          />
          <Column
            title="Done"
            status="done"
            color="border-green-500"
            emoji="✅"
            tasks={getTasksByStatus('done')}
          />
        </div>
      )}

      {/* Create Task Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-2xl max-w-lg w-full p-8 transform transition-all">
            <div className="flex justify-between items-center mb-6">
              <h3 className="text-2xl font-bold text-gray-800">Create New Task</h3>
              <button
                onClick={() => setShowCreateModal(false)}
                className="text-gray-400 hover:text-gray-600 text-2xl"
              >
                ×
              </button>
            </div>
            <form onSubmit={handleCreateTask} className="space-y-4">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Task Title
                </label>
                <input
                  type="text"
                  required
                  value={newTask.title}
                  onChange={(e) => setNewTask({ ...newTask, title: e.target.value })}
                  className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all outline-none"
                  placeholder="What needs to be done?"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Description
                </label>
                <textarea
                  value={newTask.description}
                  onChange={(e) => setNewTask({ ...newTask, description: e.target.value })}
                  className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all outline-none resize-none"
                  rows={3}
                  placeholder="Add more details..."
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Priority
                </label>
                <div className="grid grid-cols-3 gap-3">
                  {(['low', 'medium', 'high'] as const).map((priority) => (
                    <button
                      key={priority}
                      type="button"
                      onClick={() => setNewTask({ ...newTask, priority })}
                      className={`py-3 px-4 rounded-xl font-semibold transition-all ${
                        newTask.priority === priority
                          ? 'bg-gradient-to-r from-orange-300 to-orange-400 text-white shadow-lg scale-105'
                          : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                      }`}
                    >
                      {getPriorityEmoji(priority)} {priority}
                    </button>
                  ))}
                </div>
              </div>
              <div className="flex space-x-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowCreateModal(false)}
                  className="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-700 font-semibold py-3 px-4 rounded-xl transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="flex-1 bg-gradient-to-r from-orange-300 to-orange-400 hover:from-orange-400 hover:to-orange-500 text-white font-semibold py-3 px-4 rounded-xl shadow-lg transform hover:scale-105 transition-all duration-200"
                >
                  Create Task
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <style>{`
        .custom-scrollbar::-webkit-scrollbar {
          width: 6px;
        }
        .custom-scrollbar::-webkit-scrollbar-track {
          background: #f1f1f1;
          border-radius: 10px;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb {
          background: #888;
          border-radius: 10px;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb:hover {
          background: #555;
        }
      `}</style>
    </div>
  );
};

export default Board;
