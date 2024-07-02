module.exports = {
    preset: 'ts-jest',
    testEnvironment: 'node',
    transform: {
      '^.+\\.(ts|tsx)$': 'ts-jest',
      '^.+\\.jsx?$': 'babel-jest',
    },
    transformIgnorePatterns: [
      '/node_modules/(?!axios/.*)',
    ],
    moduleFileExtensions: [
      'ts',
      'tsx',
      'js',
      'jsx',
      'json',
      'node',
    ],
};
