const { Sequelize, DataTypes, Model } = require('sequelize');
const sequelize = new Sequelize('postgres://triviaUser:triviaUser@localhost:5432/trivia')


try {
    sequelize.authenticate().then(() => {
        console.log('Connection has been established successfully.');

})
  } catch (error) {
    console.error('Unable to connect to the database:', error);
  }




class Question extends Model {}

Question.init({
  // Model attributes are defined here
  firstName: {
    type: DataTypes.STRING,
    allowNull: false
  },
  lastName: {
    type: DataTypes.STRING
    // allowNull defaults to true
  }
}, {
  // Other model options go here
  sequelize, // We need to pass the connection instance
  modelName: 'User' // We need to choose the model name
});

// the defined model is the class itself
console.log(User === sequelize.models.User); // true